"use client"

import { useState } from "react"
import { Download } from "lucide-react"
import { downloadBuild } from "@/lib/api"
import { Button } from "@/components/ui/button"

interface Build {
  build_id: string
  filename: string
  timestamp: string
}

interface BuildsPanelProps {
  builds: Build[]
  loading: boolean
}

export function BuildsPanel({ builds, loading }: BuildsPanelProps) {
  const [downloading, setDownloading] = useState<string | null>(null)

  const handleDownload = async (buildId: string) => {
    setDownloading(buildId)
    try {
      await downloadBuild(buildId)
    } catch (error) {
      console.error("Download failed:", error)
    } finally {
      setDownloading(null)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-[400px]">
        <div className="animate-pulse text-green-500">Loading builds...</div>
      </div>
    )
  }

  if (builds.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-[400px]">
        <div className="text-green-500/70 mb-4">No builds found</div>
        <Button variant="outline" className="border-green-500/30 text-green-400 hover:bg-green-900/20">
          Trigger New Build
        </Button>
      </div>
    )
  }

  return (
    <div className="h-[400px] overflow-y-auto custom-scrollbar">
      <div className="grid grid-cols-1 gap-4">
        {builds.map((build) => (
          <div
            key={build.build_id}
            className="border border-green-500/30 rounded-md p-4 bg-black hover:bg-green-900/10 transition-colors"
          >
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-bold">{build.filename}</h3>
                <p className="text-sm text-green-500/70">ID: {build.build_id}</p>
                <p className="text-xs text-green-400/50 mt-1">{new Date(build.timestamp).toLocaleString()}</p>
              </div>
              <Button
                variant="outline"
                size="sm"
                className="border-green-500/30 text-green-400 hover:bg-green-900/20"
                onClick={() => handleDownload(build.build_id)}
                disabled={downloading === build.build_id}
              >
                {downloading === build.build_id ? (
                  <span className="animate-pulse">Downloading...</span>
                ) : (
                  <>
                    <Download className="h-4 w-4 mr-2" />
                    Download
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
