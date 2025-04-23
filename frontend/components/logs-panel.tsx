"use client"

import { useState } from "react"
import { Download, User, Globe, Laptop } from "lucide-react"
import { downloadLogs } from "@/lib/api"
import { Button } from "@/components/ui/button"

interface Log {
  uuid: string
  build_id: string
  timestamp: string
  path: string
  computer_name: string
  user_name: string
  country_name: string
  city_name: string
}

interface LogsPanelProps {
  logs: Log[]
  loading: boolean
}

export function LogsPanel({ logs, loading }: LogsPanelProps) {
  const [downloading, setDownloading] = useState<string | null>(null)

  const handleDownload = async (uuid: string) => {
    setDownloading(uuid)
    try {
      await downloadLogs(uuid)
    } catch (error) {
      console.error("Download failed:", error)
    } finally {
      setDownloading(null)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-[400px]">
        <div className="animate-pulse text-green-500">Loading logs...</div>
      </div>
    )
  }

  if (logs.length === 0) {
    return (
      <div className="flex items-center justify-center h-[400px]">
        <div className="text-green-500/70">No logs found</div>
      </div>
    )
  }

  return (
    <div className="h-[400px] overflow-y-auto custom-scrollbar">
      <div className="grid grid-cols-1 gap-4">
        {logs.map((log) => (
          <div
            key={log.uuid}
            className="border border-green-500/30 rounded-md p-4 bg-black hover:bg-green-900/10 transition-colors"
          >
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-bold flex items-center">
                  <Laptop className="h-4 w-4 mr-2" />
                  {log.computer_name}
                </h3>
                <p className="text-sm text-green-500/70 mt-1">UUID: {log.uuid}</p>
                <p className="text-sm text-green-400/70 mt-1">Build ID: {log.build_id}</p>
                <div className="flex items-center mt-2 text-xs text-green-400/50">
                  <User className="h-3 w-3 mr-1" />
                  <span className="mr-3">{log.user_name}</span>
                  <Globe className="h-3 w-3 mr-1" />
                  <span>
                    {log.country_name}, {log.city_name}
                  </span>
                </div>
                <p className="text-xs text-green-400/50 mt-1">{new Date(log.timestamp).toLocaleString()}</p>
              </div>
              <Button
                variant="outline"
                size="sm"
                className="border-green-500/30 text-green-400 hover:bg-green-900/20"
                onClick={() => handleDownload(log.uuid)}
                disabled={downloading === log.uuid}
              >
                {downloading === log.uuid ? (
                  <span className="animate-pulse">Downloading...</span>
                ) : (
                  <>
                    <Download className="h-4 w-4 mr-2" />
                    Download Logs
                  </>
                )}
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
