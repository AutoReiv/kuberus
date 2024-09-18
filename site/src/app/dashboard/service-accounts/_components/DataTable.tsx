import React, { useState } from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { MoreHorizontal, FileText, Copy, UserPlus, UserMinus } from "lucide-react"
import Link from "next/link"
import { motion } from "framer-motion"

const DataTable = ({ serviceAccounts }) => {
  const [viewMode, setViewMode] = useState("table")
  const [filter, setFilter] = useState("")

  const filteredServiceAccounts = serviceAccounts.filter(sa =>
    sa.metadata.name.toLowerCase().includes(filter.toLowerCase())
  )

  const TableView = () => (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Namespace</TableHead>
          <TableHead>Created At</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {filteredServiceAccounts.map((sa) => (
          <TableRow key={sa.metadata.uid}>
            <TableCell>{sa.metadata.name}</TableCell>
            <TableCell>{sa.metadata.namespace}</TableCell>
            <TableCell>{new Date(sa.metadata.creationTimestamp).toLocaleString()}</TableCell>
            <TableCell>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="h-8 w-8 p-0">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuLabel>Actions</DropdownMenuLabel>
                  <DropdownMenuItem>
                    <FileText className="mr-2 h-4 w-4" />
                    <Link href={`/dashboard/service-accounts/${sa.metadata.namespace}/${sa.metadata.name}`}>
                      View
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => console.log("Clone", sa.metadata.name)}>
                    <Copy className="mr-2 h-4 w-4" />
                    Clone
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => console.log("Add Token", sa.metadata.name)}>
                    <UserPlus className="mr-2 h-4 w-4" />
                    Add Token
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => console.log("Remove Token", sa.metadata.name)}>
                    <UserMinus className="mr-2 h-4 w-4" />
                    Remove Token
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )

  const GridView = () => (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {filteredServiceAccounts.map((sa) => (
        <motion.div
          key={sa.metadata.uid}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <Card>
            <CardHeader className="bg-primary/10 rounded-t-lg">
              <CardTitle className="text-xl font-bold">
                {sa.metadata.name}
              </CardTitle>
              <CardDescription className="text-sm opacity-70">
                {sa.metadata.namespace}
              </CardDescription>
            </CardHeader>
            <CardContent className="flex-grow">
              <p className="text-sm mt-2">
                Created: {new Date(sa.metadata.creationTimestamp).toLocaleString()}
              </p>
            </CardContent>
            <CardContent className="bg-secondary/10 rounded-b-lg">
              <Link
                href={`/dashboard/service-accounts/${sa.metadata.namespace}/${sa.metadata.name}`}
                className="w-full"
              >
                <Button className="w-full">View Details</Button>
              </Link>
            </CardContent>
          </Card>
        </motion.div>
      ))}
    </div>
  )

  return (
    <Card className="h-full">
      <CardHeader>
        <div className="justify-between item-start flex">
          <div className="flex flex-col gap-4">
            <CardTitle className="font-bold">Service Accounts</CardTitle>
            <CardDescription>Manage Service Accounts</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between py-4">
          <Input
            placeholder="Filter service accounts..."
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="max-w-sm"
          />
          <Button
            onClick={() => setViewMode(viewMode === "table" ? "grid" : "table")}
          >
            {viewMode === "table" ? "Switch to Grid" : "Switch to Table"}
          </Button>
        </div>
        {viewMode === "table" ? <TableView /> : <GridView />}
      </CardContent>
    </Card>
  )
}

export default DataTable
