"use client"

import { useEffect, useRef } from "react"

interface TerminalProps {
  lines: string[]
}

export function Terminal({ lines }: TerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight
    }
  }, [lines])

  
  const formatLine = (line: string) => {
    if (line.startsWith(">")) {
      return (
        <div className="flex">
          <span className="text-green-500 mr-2">$</span>
          <span>{line.substring(2)}</span>
        </div>
      )
    } else if (line.includes("downloaded successfully")) {
      return <div className="text-green-400">{line}</div>
    } else if (line.includes("Error:") || line.includes("failed")) {
      return <div className="text-red-400">{line}</div>
    } else if (line.includes("UUID:") || line.includes("ID:")) {
      return (
        <div>
          {line.split("|").map((segment, i) => (
            <span key={i} className={i === 0 ? "text-yellow-300" : ""}>
              {segment}
              {i < line.split("|").length - 1 ? "|" : ""}
            </span>
          ))}
        </div>
      )
    } else if (line.startsWith("Available commands:") || line.startsWith("  ")) {
      return <div className="text-blue-300">{line}</div>
    } else if (line.includes("Downloading")) {
      return <div className="text-yellow-300">{line}</div>
    } else {
      return <div>{line}</div>
    }
  }

  return (
    <div
      ref={terminalRef}
      className="bg-black text-green-400 p-4 font-mono text-sm h-[400px] overflow-y-auto custom-scrollbar"
    >
      {lines.map((line, index) => (
        <div key={index} className="mb-1">
          {formatLine(line)}
        </div>
      ))}
    </div>
  )
}
