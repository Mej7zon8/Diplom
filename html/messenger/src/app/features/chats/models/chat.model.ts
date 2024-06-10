import {Message} from "./message.model";

export type ChatList = Chat[] | null

export interface Chat {
  id: string
  is_group: boolean
  name: string
  participants: Participant[]
  last_message: LastMessage
  unread_count: number
}

export interface Participant {
  id: string
  name: string
}

export type LastMessage = Message
