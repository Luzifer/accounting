export type AccountType = 'budget' | 'category' | 'tracking'

export interface Account {
  balance: number
  hidden: boolean
  id: string
  name: string
  type: AccountType
}

export interface DateRange {
  end?: Date
  start?: Date
}

export interface DateRangeStorage {
  end: string
  start: string
}

export interface Transaction {
  account: string | null
  amount: number
  category: string | null
  cleared: boolean
  description: string
  id: string
  payee: string
  reconciled: boolean
  time: string
}

export interface JsonPatchOperation {
  op: 'replace'
  path: string
  value: unknown
}
