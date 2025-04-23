"use client"

import { useEffect, useRef, useState } from "react"

interface Log {
  uuid: string
  country_name: string
  city_name: string
  timestamp: string
}

interface GeoMapProps {
  logs: Log[]
}

interface MapPoint {
  x: number
  y: number
  country: string
  city: string
  timestamp: string
  uuid: string
}

export function GeoMap({ logs }: GeoMapProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [selectedPoint, setSelectedPoint] = useState<MapPoint | null>(null)
  const [points, setPoints] = useState<MapPoint[]>([])
  const [hoveredPoint, setHoveredPoint] = useState<MapPoint | null>(null)

  
  useEffect(() => {
    
    const newPoints = logs.map((log) => {
      
      
      const hash = Array.from(log.uuid).reduce((acc, char) => acc + char.charCodeAt(0), 0)
      const x = ((hash % 100) / 100) * 0.8 + 0.1 
      const y = (((hash * 13) % 100) / 100) * 0.8 + 0.1 

      return {
        x,
        y,
        country: log.country_name,
        city: log.city_name,
        timestamp: log.timestamp,
        uuid: log.uuid,
      }
    })

    setPoints(newPoints)
  }, [logs])

  
  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext("2d")
    if (!ctx) return

    
    const resizeCanvas = () => {
      const parent = canvas.parentElement
      if (parent) {
        canvas.width = parent.clientWidth
        canvas.height = parent.clientHeight
      }
    }

    resizeCanvas()
    window.addEventListener("resize", resizeCanvas)

    
    const drawMap = () => {
      if (!ctx || !canvas) return

      
      ctx.clearRect(0, 0, canvas.width, canvas.height)

      
      ctx.strokeStyle = "rgba(76, 175, 80, 0.1)"
      ctx.lineWidth = 1

      
      for (let i = 0; i <= 10; i++) {
        const x = (canvas.width / 10) * i
        ctx.beginPath()
        ctx.moveTo(x, 0)
        ctx.lineTo(x, canvas.height)
        ctx.stroke()
      }

      
      for (let i = 0; i <= 10; i++) {
        const y = (canvas.height / 10) * i
        ctx.beginPath()
        ctx.moveTo(0, y)
        ctx.lineTo(canvas.width, y)
        ctx.stroke()
      }

      
      ctx.strokeStyle = "rgba(76, 175, 80, 0.3)"
      ctx.lineWidth = 1

      for (let i = 0; i < points.length; i++) {
        for (let j = i + 1; j < points.length; j++) {
          if (i !== j) {
            ctx.beginPath()
            ctx.moveTo(points[i].x * canvas.width, points[i].y * canvas.height)
            ctx.lineTo(points[j].x * canvas.width, points[j].y * canvas.height)
            ctx.stroke()
          }
        }
      }

      
      points.forEach((point, index) => {
        const x = point.x * canvas.width
        const y = point.y * canvas.height

        
        ctx.beginPath()
        ctx.arc(x, y, 10, 0, Math.PI * 2)
        ctx.fillStyle = "rgba(76, 175, 80, 0.2)"
        ctx.fill()

        
        ctx.beginPath()
        ctx.arc(x, y, 5, 0, Math.PI * 2)

        if (hoveredPoint === point || selectedPoint === point) {
          ctx.fillStyle = "#4caf50"
        } else {
          ctx.fillStyle = "rgba(76, 175, 80, 0.8)"
        }

        ctx.fill()

        
        if (selectedPoint === point) {
          ctx.beginPath()
          const pulseSize = 15 + Math.sin(Date.now() / 200) * 5
          ctx.arc(x, y, pulseSize, 0, Math.PI * 2)
          ctx.strokeStyle = "rgba(76, 175, 80, 0.5)"
          ctx.stroke()
        }
      })
    }

    
    let animationId: number
    const animate = () => {
      drawMap()
      animationId = requestAnimationFrame(animate)
    }

    animate()

    
    const handleMouseMove = (e: MouseEvent) => {
      const rect = canvas.getBoundingClientRect()
      const mouseX = (e.clientX - rect.left) / canvas.width
      const mouseY = (e.clientY - rect.top) / canvas.height

      
      let found = false
      for (const point of points) {
        const dx = point.x - mouseX
        const dy = point.y - mouseY
        const distance = Math.sqrt(dx * dx + dy * dy)

        if (distance < 0.02) {
          
          setHoveredPoint(point)
          found = true
          break
        }
      }

      if (!found) {
        setHoveredPoint(null)
      }
    }

    const handleClick = (e: MouseEvent) => {
      if (hoveredPoint) {
        setSelectedPoint(hoveredPoint)
      } else {
        setSelectedPoint(null)
      }
    }

    canvas.addEventListener("mousemove", handleMouseMove)
    canvas.addEventListener("click", handleClick)

    return () => {
      window.removeEventListener("resize", resizeCanvas)
      canvas.removeEventListener("mousemove", handleMouseMove)
      canvas.removeEventListener("click", handleClick)
      cancelAnimationFrame(animationId)
    }
  }, [points, hoveredPoint, selectedPoint])

  return (
    <div className="relative h-full w-full">
      <canvas ref={canvasRef} className="w-full h-full cursor-crosshair" />

      {hoveredPoint && (
        <div
          className="absolute bg-black border border-green-500/30 p-2 rounded-md text-xs pointer-events-none"
          style={{
            left: `${hoveredPoint.x * 100}%`,
            top: `${hoveredPoint.y * 100}%`,
            transform: "translate(-50%, -130%)",
          }}
        >
          <div className="font-bold">
            {hoveredPoint.city}, {hoveredPoint.country}
          </div>
        </div>
      )}

      {selectedPoint && (
        <div
          className="absolute bg-black border border-green-500/30 p-3 rounded-md text-sm z-10"
          style={{
            left: 20,
            bottom: 20,
            maxWidth: "300px",
          }}
        >
          <div className="font-bold text-green-400 mb-1">
            {selectedPoint.city}, {selectedPoint.country}
          </div>
          <div className="text-xs text-green-500/70 mb-1">UUID: {selectedPoint.uuid}</div>
          <div className="text-xs text-green-400/50">{new Date(selectedPoint.timestamp).toLocaleString()}</div>
          <button
            className="absolute top-2 right-2 text-green-500/70 hover:text-green-500"
            onClick={() => setSelectedPoint(null)}
          >
            ×
          </button>
        </div>
      )}

      <div className="absolute bottom-2 right-2 text-xs text-green-500/50">
        {points.length} locations • Click on a point for details
      </div>
    </div>
  )
}
