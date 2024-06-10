export type Messages = Message[] | null

export interface Message {
  id: number
  created: Date | string
  sender: string
  content: MessageContent
}

export interface MessageContent {
  text: string
}
