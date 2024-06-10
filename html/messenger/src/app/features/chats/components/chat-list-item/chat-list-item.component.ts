import {Component, EventEmitter, Input, Output} from '@angular/core';
import {Chat} from "../../models/chat.model";
import {uuid2Colors} from "../../../../core/tools/colors";
import {formatDate} from "@angular/common";

@Component({
  selector: 'app-chat-list-item',
  standalone: true,
  imports: [],
  templateUrl: './chat-list-item.component.html',
  styleUrl: './chat-list-item.component.scss'
})
export class ChatListItemComponent {
  protected model?: {
    chat: Chat
    avatar: {
      letter: string
      background: string
      foreground: string
    }
  }

  @Input() selected = false

  @Input()
  set data(v: Chat) {
    const {bg, fg} = uuid2Colors(v.id)

    this.model = {
      chat: v,
      avatar: {
        letter: v.name ? v.name[0] : '',
        background: bg,
        foreground: fg
      }
    }
  }

  get lastMessageDate() {
    const lastMessageDateString = this.model?.chat.last_message.created
    // If the date is not valid, return an empty string.
    if (!lastMessageDateString || lastMessageDateString == "0001-01-01T00:00:00Z") {
      return ''
    }

    // If the last message was sent less than 24 hours ago, show the time.
    const now = new Date()
    const lastMessageDate = new Date(lastMessageDateString)
    if (now.getTime() - lastMessageDate.getTime() < 24 * 60 * 60 * 1000) {
      return formatDate(lastMessageDate, 'HH:mm', 'en-US')
    }
    // Otherwise, show the date.
    return formatDate(lastMessageDate, 'd MMM', 'en-US')
  }
}
