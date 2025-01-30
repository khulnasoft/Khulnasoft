import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { useState } from 'react'
import { useEffect } from 'react'
import { Team } from '@/utils/useUser'
import { getAPIUrl } from '@/app/(dashboard)/dashboard/utils'

interface Sandbox {
  alias: string
  clientID: string
  cpuCount: number
  endAt: string
  memoryMB: number
  metadata: Record<string, any>
  sandboxID: string
  startedAt: string
  templateID: string
}

export function SandboxesContent({
  team,
  domain,
}: {
  team: Team
  domain: string
}) {
  const [runningSandboxes, setRunningSandboxes] = useState<Sandbox[]>([])

  useEffect(() => {
    function f() {
      const apiKey = team.apiKeys[0]
      if (apiKey) {
        fetchSandboxes(domain, apiKey).then((newSandboxes) => {
          if (newSandboxes) {
            setRunningSandboxes(newSandboxes)
          }
        })
      }
    }

    const interval = setInterval(() => {
      f()
    }, 5000)

    f()
    // Cleanup interval on component unmount
    return () => clearInterval(interval)
  }, [team])

  return (
    <div className="flex flex-col justify-center">
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-orange-500/10 dark:hover:bg-orange-500/10 border-b border-white/5 ">
            <TableHead>Sandbox ID</TableHead>
            <TableHead>Template ID</TableHead>
            <TableHead>Alias</TableHead>
            <TableHead>Started at</TableHead>
            <TableHead>End at</TableHead>
            <TableHead>vCPUs</TableHead>
            <TableHead>RAM MiB</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {runningSandboxes.length === 0 ? (
            <TableRow>
              <TableCell colSpan={8} className="text-center">
                No running sandboxes
              </TableCell>
            </TableRow>
          ) : (
            runningSandboxes.map((sandbox) => (
              <TableRow
                className="hover:bg-orange-300/10 dark:hover:bg-orange-300/10 border-b border-white/5"
                key={sandbox.sandboxID}
              >
                <TableCell>{sandbox.sandboxID}</TableCell>
                <TableCell>{sandbox.templateID}</TableCell>
                <TableCell>{sandbox.alias}</TableCell>
                <TableCell>
                  {new Date(sandbox.startedAt).toLocaleString()}
                </TableCell>
                <TableCell>
                  {new Date(sandbox.endAt).toLocaleString()}
                </TableCell>
                <TableCell>{sandbox.cpuCount}</TableCell>
                <TableCell>{sandbox.memoryMB}</TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}

async function fetchSandboxes(
  domain: string,
  apiKey: string
): Promise<Sandbox[]> {
  const res = await fetch(getAPIUrl(domain, '/sandboxes'), {
    method: 'GET',
    headers: {
      'X-API-KEY': apiKey,
    },
  })
  try {
    const data: Sandbox[] = await res.json()

    // Latest sandboxes first
    return data.sort(
      (a, b) =>
        new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime()
    )
  } catch (e) {
    // TODO: add sentry event here
    return []
  }
}
