import {Component, EventEmitter, Input, Output} from '@angular/core';
import {Chat, ChatList} from "../../models/chat.model";
import {MatIconButton} from "@angular/material/button";
import {MatIcon} from "@angular/material/icon";
import {ChatListItemComponent} from "../chat-list-item/chat-list-item.component";
import {MatFormField} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {ContextMenuComponent, contextMenuItem} from "../../../../shared/components/context-menu/context-menu.component";

@Component({
  selector: 'app-chat-list',
  standalone: true,
  imports: [
    MatIconButton,
    MatIcon,
    ChatListItemComponent,
    MatFormField,
    MatInput,
    ContextMenuComponent
  ],
  templateUrl: './chat-list.component.html',
  styleUrl: './chat-list.component.scss'
})
export class ChatListComponent {
  @Input() data: ChatList = null
  @Output() chatSelected = new EventEmitter<Chat>()
  @Output() chatLeave = new EventEmitter<Chat>()
  @Output() createChatClick = new EventEmitter<void>()

  protected model = {
    selected: "",
    search: ""
  }

  protected satisfies(chat: Chat) {
    return this.model.search == "" || chat.name.toLowerCase().includes(this.model.search.toLowerCase())
  }

  protected generateContextMenu(chat: Chat): contextMenuItem[] {
    return [{
      text: "Выйти из чата",
      disabled: !this.canLeaveChat(chat),
      action: () => this.leaveChat(chat)
    }]
  }

  protected leaveChat(chat: Chat) {
    this.chatLeave.emit(chat)
  }

  protected canLeaveChat(chat: Chat) {
    return chat.is_group
  }
}
