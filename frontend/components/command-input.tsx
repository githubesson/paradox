"use client"

import { useState, type KeyboardEvent } from "react"

interface CommandInputProps {
  onSubmit: (command: string) => void
}

export function CommandInput({ onSubmit }: CommandInputProps) {
  const [command, setCommand] = useState("")
  const [commandHistory, setCommandHistory] = useState<string[]>([])
  const [historyIndex, setHistoryIndex] = useState(-1)

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && command.trim()) {
      onSubmit(command)
      setCommandHistory((prev) => [command, ...prev])
      setCommand("")
      setHistoryIndex(-1)
    } else if (e.key === "ArrowUp") {
      e.preventDefault()
      if (commandHistory.length > 0 && historyIndex < commandHistory.length - 1) {
        const newIndex = historyIndex + 1
        setHistoryIndex(newIndex)
        setCommand(commandHistory[newIndex])
      }
    } else if (e.key === "ArrowDown") {
      e.preventDefault()
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1
        setHistoryIndex(newIndex)
        setCommand(commandHistory[newIndex])
      } else if (historyIndex === 0) {
        setHistoryIndex(-1)
        setCommand("")
      }
    } else if (e.key === "Tab") {
      e.preventDefault()

      
      const commandStart = command.toLowerCase().trim()

      if (commandStart === "") return

      const completions = [
        "help",
        "build",
        "logs",
        "builds",
        "clear",
        "download build ",
        "download log ",
        "list builds",
        "list logs",
      ]

      const matches = completions.filter((c) => c.startsWith(commandStart))

      if (matches.length === 1) {
        setCommand(matches[0])
      } else if (matches.length > 1 && commandStart === "download") {
        setCommand("download ")
      } else if (matches.length > 1 && commandStart === "list") {
        setCommand("list ")
      }
    }
  }

  return (
    <div className="flex items-center p-2 border-t border-green-500/30 bg-black">
      <span className="text-green-500 mr-2">$</span>
      <input
        type="text"
        value={command}
        onChange={(e) => setCommand(e.target.value)}
        onKeyDown={handleKeyDown}
        className="flex-1 bg-transparent border-none outline-none text-green-400 font-mono"
        placeholder="Enter command..."
        autoFocus
      />
    </div>
  )
}
