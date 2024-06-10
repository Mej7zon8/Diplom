import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {ApiResponse} from "../../../core/models/api-generic.model";
import {resolve} from "../../../environments/environment";
import {Chat, ChatList} from "../models/chat.model";

@Injectable({
  providedIn: 'root'
})
export class ChatService {

  constructor(private http: HttpClient) { }

  list() {
    const fd = this.getFormData("list")
    return this.http.post<ApiResponse<ChatList>>(resolve(), fd)
  }

  create(participants: string[], name: string) {
    const fd = this.getFormData("create")
    fd.set("data", JSON.stringify({"participants": participants, "name": name}))
    return this.http.post<ApiResponse<Chat>>(resolve(), fd)
  }

  leave(chat: Chat) {
    const fd = this.getFormData("leave")
    fd.set("data", JSON.stringify({"chat": chat.id}))
    return this.http.post<ApiResponse<Chat>>(resolve(), fd)
  }

  private getFormData(method: string) {
    let fd = new FormData()
    fd.set("module", "chats")
    fd.set("method", method)
    return fd
  }
}
