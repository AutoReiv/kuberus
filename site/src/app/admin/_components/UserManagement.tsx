import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import React from 'react'

const UserManagement = () => {
  return (
    <Card>
      <CardHeader>Create User</CardHeader>
      <CardContent>
        <Button>Create</Button>
      </CardContent>
    </Card>
  )
}

export default UserManagement