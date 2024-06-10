import {Component, EventEmitter, Input} from '@angular/core';
import {Chat} from "../../models/chat.model";
import {Message, Messages} from "../../models/message.model";
import {MessageService} from "../../services/message.service";
import {SnackbarErrorHandlerService} from "../../../../core/services/snackbar-error-handler.service";
import {exhaustMap, merge, Observable, of, tap, timer} from "rxjs";
import {ApiResponse} from "../../../../core/models/api-generic.model";
import {MatFormField, MatSuffix} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {MatIconButton} from "@angular/material/button";
import {MatIcon} from "@angular/material/icon";
import {AsyncPipe} from "@angular/common";
import {ChatMessageComponent} from "../chat-message/chat-message.component";
import {AuthenticationService} from "../../../authentication/services/authentication.service";

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [
    MatFormField,
    MatInput,
    MatIconButton,
    MatSuffix,
    MatIcon,
    AsyncPipe,
    ChatMessageComponent
  ],
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.scss'
})
export class ChatComponent {
  protected get disabled() {
    return this.model.chat == null || this.model.state == "sending"
  }

  protected manualUpdates = new EventEmitter<boolean>()
  protected dataEmitter: Observable<ApiResponse<Messages>>
  protected model: {
    chat: Chat | null,
    isGroup: boolean,
    state: "sending" | "idle"
    username2name: {[key: string]: string}
  } = {
    chat: null,
    isGroup: false,
    state: "idle",
    username2name: {}
  }

  constructor(private service: MessageService, private eh: SnackbarErrorHandlerService, private authService: AuthenticationService) {
    this.dataEmitter = merge(timer(0, 5000), this.manualUpdates).pipe(
      exhaustMap(src => {
        const emitEmpty = typeof src == "boolean" && src
        if (emitEmpty || this.model.chat == null) {
          const emptyResponse: ApiResponse<Messages> = {
            data: [],
            error: "",
            status: 200
          }
          return of(emptyResponse)
        }
        return this.service.read(this.model.chat.id, -1, -1)
      }),
      tap(this.eh.tap())
    )
  }

  @Input()
  set chat(chat: Chat | null) {
    this.model.chat = chat
    this.model.isGroup = chat != null && chat.participants.length > 1
    this.model.username2name = {}
    chat?.participants.forEach(p => this.model.username2name[p.id] = p.name)
    this.manualUpdates.emit(chat == null || this.model.chat == null || chat.id != this.model.chat.id)
  }

  protected sendMessage(el: HTMLInputElement) {
    this.model.state = "sending"
    this.service.send(this.model.chat!.id, {text: el.value})
      .pipe(tap(this.eh.tap()))
      .subscribe(v => {
        this.model.state = "idle"
        if (v.error == "") {
          this.manualUpdates.emit()
          el.value = ""
          setTimeout(() => el.focus(), 0)
        }
      })
  }

  protected getUserName(message: Message) {
    if (message.sender == this.authService.getUsername())
      return ""
    return this.model.username2name[message.sender] || message.sender
  }
}
