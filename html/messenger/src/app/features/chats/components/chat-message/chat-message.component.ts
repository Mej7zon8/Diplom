import {Component, Input} from '@angular/core';
import {Message} from "../../models/message.model";
import {DatePipe, JsonPipe} from "@angular/common";
import {AuthenticationService} from "../../../authentication/services/authentication.service";
import {MatIcon} from "@angular/material/icon";

@Component({
  selector: 'app-chat-message',
  standalone: true,
  imports: [
    JsonPipe,
    DatePipe,
    MatIcon
  ],
  templateUrl: './chat-message.component.html',
  styleUrl: './chat-message.component.scss'
})
export class ChatMessageComponent {
  @Input() data!: Message
  @Input() name: string = ""

  protected model = {
    us: ""
  }

  constructor(authService: AuthenticationService) {
    this.model.us = authService.getUsername() || ""
  }

  get sentByUs() {
    return this.data.sender == this.model.us
  }
}
