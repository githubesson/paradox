"use client"

import { useEffect, useState, useRef } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth"
import { Terminal } from "@/components/terminal"
import { BuildsPanel } from "@/components/builds-panel"
import { LogsPanel } from "@/components/logs-panel"
import { CommandInput } from "@/components/command-input"
import { GeoMap } from "@/components/geo-map"
import { fetchBuilds, fetchLogs, triggerBuild, downloadBuild as apiDownloadBuild, downloadLogs as apiDownloadLogs, Log, Build } from "@/lib/api"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Cpu, Database, Globe, LogOut, TerminalIcon } from "lucide-react"

export default function Dashboard() {
  const router = useRouter()
  const { user, logout, isLoading: authLoading } = useAuth()
  const [logs, setLogs] = useState<Log[]>([])
  const [builds, setBuilds] = useState<Build[]>([])
  const [loading, setLoading] = useState(true)
  const hasInitialized = useRef(false)
  const [terminalOutput, setTerminalOutput] = useState<string[]>([
    "Paradox Terminal v1.0.0",
    "Type 'help' for available commands",
    "Initializing system...",
  ])

  useEffect(() => {
    
    if (!authLoading && !user) {
      router.push('/login')
      return
    }

    
    if (user && !hasInitialized.current) {
      const loadData = async () => {
        try {
          const [logsData, buildsData] = await Promise.all([fetchLogs(), fetchBuilds()])
          setLogs(logsData.logs || [])
          setBuilds(buildsData.builds || [])
          setTerminalOutput((prev) => [...prev, "System initialized. Ready for commands."])
          hasInitialized.current = true
        } catch (error: unknown) {
          const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
          setTerminalOutput((prev) => [...prev, `Error: ${errorMessage}`])
        } finally {
          setLoading(false)
        }
      }

      loadData()
    }
  }, [user, authLoading, router])

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  const handleCommand = async (command: string) => {
    setTerminalOutput((prev) => [...prev, `> ${command}`])

    const cmd = command.toLowerCase().trim()

    if (cmd === "help") {
      setTerminalOutput((prev) => [
        ...prev,
        "Available commands:",
        "  help - Show this help message",
        "  build - Trigger a new build",
        "  logs - Show recent logs",
        "  builds - Show recent builds",
        "  list builds - List all available builds with details",
        "  list logs - List all available logs with details",
        "  download build <build_id> - Download a specific build",
        "  download log <uuid> - Download logs for a specific UUID",
        "  clear - Clear terminal",
      ])
    } else if (cmd === "build") {
      setTerminalOutput((prev) => [...prev, "Triggering new build..."])
      try {
        const result = await triggerBuild()
        setTerminalOutput((prev) => [...prev, `Build successful: ${result.build_id}`])
        
        const buildsData = await fetchBuilds()
        setBuilds(buildsData.builds || [])
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Build failed: ${errorMessage}`])
      }
    } else if (cmd === "logs") {
      setTerminalOutput((prev) => [...prev, "Fetching recent logs..."])
      try {
        const logsData = await fetchLogs()
        setLogs(logsData.logs || [])
        setTerminalOutput((prev) => [...prev, `Retrieved ${logsData.count} logs`])
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Failed to fetch logs: ${errorMessage}`])
      }
    } else if (cmd === "builds") {
      setTerminalOutput((prev) => [...prev, "Fetching recent builds..."])
      try {
        const buildsData = await fetchBuilds()
        setBuilds(buildsData.builds || [])
        setTerminalOutput((prev) => [...prev, `Retrieved ${buildsData.count} builds`])
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Failed to fetch builds: ${errorMessage}`])
      }
    } else if (cmd === "clear") {
      setTerminalOutput(["Terminal cleared"])
    } else if (cmd.startsWith("download build ")) {
      const buildId = command.substring("download build ".length).trim()
      if (!buildId) {
        setTerminalOutput((prev) => [...prev, "Error: Please provide a build ID"])
        return
      }
      try {
        setTerminalOutput((prev) => [...prev, `Downloading build ${buildId}...`])
        await apiDownloadBuild(buildId)
        setTerminalOutput((prev) => [...prev, `Build ${buildId} downloaded successfully`])
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Error downloading build: ${errorMessage}`])
      }
    } else if (cmd.startsWith("download log ")) {
      const uuid = command.substring("download log ".length).trim()
      if (!uuid) {
        setTerminalOutput((prev) => [...prev, "Error: Please provide a log UUID"])
        return
      }
      try {
        setTerminalOutput((prev) => [...prev, `Downloading logs for ${uuid}...`])
        await apiDownloadLogs(uuid)
        setTerminalOutput((prev) => [...prev, `Logs for ${uuid} downloaded successfully`])
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Error downloading logs: ${errorMessage}`])
      }
    } else if (cmd === "list builds") {
      setTerminalOutput((prev) => [...prev, "Listing available builds:"])
      try {
        const buildsData = await fetchBuilds()
        if (buildsData.builds && buildsData.builds.length > 0) {
          buildsData.builds.forEach((build: Build, index: number) => {
            setTerminalOutput((prev) => [
              ...prev,
              `${index + 1}. ID: ${build.build_id} | File: ${build.filename} | Date: ${new Date(build.timestamp).toLocaleString()}`,
            ])
          })
        } else {
          setTerminalOutput((prev) => [...prev, "No builds available"])
        }
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Failed to fetch builds: ${errorMessage}`])
      }
    } else if (cmd === "list logs") {
      setTerminalOutput((prev) => [...prev, "Listing available logs:"])
      try {
        const logsData = await fetchLogs()
        if (logsData.logs && logsData.logs.length > 0) {
          logsData.logs.forEach((log: Log, index: number) => {
            setTerminalOutput((prev) => [
              ...prev,
              `${index + 1}. UUID: ${log.uuid} | Computer: ${log.computer_name} | Location: ${log.country_name}, ${log.city_name}`,
            ])
          })
        } else {
          setTerminalOutput((prev) => [...prev, "No logs available"])
        }
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'An unknown error occurred'
        setTerminalOutput((prev) => [...prev, `Failed to fetch logs: ${errorMessage}`])
      }
    } else {
      setTerminalOutput((prev) => [...prev, `Unknown command: ${command}`])
    }
  }

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-black text-green-400">
        <div className="text-center">
          <div className="animate-pulse">Loading...</div>
        </div>
      </div>
    )
  }

  if (!user) {
    return null 
  }

  return (
    <div className="min-h-screen bg-black text-green-400 font-mono p-4">
      <header className="border-b border-green-500/30 pb-4 mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-4xl font-bold tracking-tighter glitch-text">PARADOX TERMINAL</h1>
            <p className="text-green-500/70">Secure Server Management Interface</p>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-green-500/70">Logged in as {user.username}</span>
            <Button
              variant="outline"
              className="border-green-500/30 hover:bg-green-900/20"
              onClick={handleLogout}
            >
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </div>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Tabs defaultValue="terminal" className="w-full">
            <TabsList className="grid grid-cols-4 bg-black border border-green-500/30">
              <TabsTrigger value="terminal" className="data-[state=active]:bg-green-900/20">
                <TerminalIcon className="mr-2 h-4 w-4" />
                Terminal
              </TabsTrigger>
              <TabsTrigger value="builds" className="data-[state=active]:bg-green-900/20">
                <Cpu className="mr-2 h-4 w-4" />
                Builds
              </TabsTrigger>
              <TabsTrigger value="logs" className="data-[state=active]:bg-green-900/20">
                <Database className="mr-2 h-4 w-4" />
                Logs
              </TabsTrigger>
              <TabsTrigger value="map" className="data-[state=active]:bg-green-900/20">
                <Globe className="mr-2 h-4 w-4" />
                Map
              </TabsTrigger>
            </TabsList>

            <TabsContent value="terminal" className="border border-green-500/30 rounded-md mt-2">
              <Terminal lines={terminalOutput} />
              <CommandInput onSubmit={handleCommand} />
            </TabsContent>

            <TabsContent value="builds" className="border border-green-500/30 rounded-md mt-2 p-4">
              <BuildsPanel builds={builds} loading={loading} />
            </TabsContent>

            <TabsContent value="logs" className="border border-green-500/30 rounded-md mt-2 p-4">
              <LogsPanel logs={logs} loading={loading} />
            </TabsContent>

            <TabsContent value="map" className="border border-green-500/30 rounded-md mt-2 p-4 h-[500px]">
              <GeoMap logs={logs} />
            </TabsContent>
          </Tabs>
        </div>

        <div className="lg:col-span-1 space-y-6">
          <div className="border border-green-500/30 rounded-md p-4">
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <Cpu className="mr-2 h-5 w-5" />
              Recent Builds
            </h2>
            <div className="space-y-2 max-h-[200px] overflow-y-auto custom-scrollbar">
              {builds.slice(0, 5).map((build, index) => (
                <div
                  key={index}
                  className="p-2 border border-green-500/20 rounded-md bg-black hover:bg-green-900/10 transition-colors"
                >
                  <div className="flex justify-between">
                    <span className="text-xs">{build.build_id}</span>
                    <span className="text-xs text-green-300">{new Date(build.timestamp).toLocaleString()}</span>
                  </div>
                  <div className="text-sm mt-1">{build.filename}</div>
                </div>
              ))}
              {builds.length === 0 && !loading && (
                <div className="text-center py-4 text-green-500/50">No builds found</div>
              )}
              {loading && <div className="text-center py-4 text-green-500/50">Loading builds...</div>}
            </div>
          </div>

          <div className="border border-green-500/30 rounded-md p-4">
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <Database className="mr-2 h-5 w-5" />
              Recent Logs
            </h2>
            <div className="space-y-2 max-h-[200px] overflow-y-auto custom-scrollbar">
              {logs.slice(0, 5).map((log, index) => (
                <div
                  key={index}
                  className="p-2 border border-green-500/20 rounded-md bg-black hover:bg-green-900/10 transition-colors"
                >
                  <div className="flex justify-between">
                    <span className="text-xs truncate max-w-[150px]">{log.uuid}</span>
                    <span className="text-xs text-green-300">{new Date(log.timestamp).toLocaleString()}</span>
                  </div>
                  <div className="text-sm mt-1 flex justify-between">
                    <span>{log.computer_name}</span>
                    <span className="text-green-300/70">{log.country_name}</span>
                  </div>
                </div>
              ))}
              {logs.length === 0 && !loading && <div className="text-center py-4 text-green-500/50">No logs found</div>}
              {loading && <div className="text-center py-4 text-green-500/50">Loading logs...</div>}
            </div>
          </div>

          <div className="border border-green-500/30 rounded-md p-4">
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <TerminalIcon className="mr-2 h-5 w-5" />
              Quick Actions
            </h2>
            <div className="grid grid-cols-2 gap-2">
              <button
                onClick={() => handleCommand("build")}
                className="p-2 border border-green-500/30 rounded-md hover:bg-green-900/20 transition-colors"
              >
                New Build
              </button>
              <button
                onClick={() => handleCommand("logs")}
                className="p-2 border border-green-500/30 rounded-md hover:bg-green-900/20 transition-colors"
              >
                Refresh Logs
              </button>
              <button
                onClick={() => handleCommand("builds")}
                className="p-2 border border-green-500/30 rounded-md hover:bg-green-900/20 transition-colors"
              >
                Refresh Builds
              </button>
              <button
                onClick={() => handleCommand("clear")}
                className="p-2 border border-green-500/30 rounded-md hover:bg-green-900/20 transition-colors"
              >
                Clear Terminal
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
