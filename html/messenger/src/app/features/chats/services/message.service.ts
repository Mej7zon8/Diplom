import { Injectable } from '@angular/core';
import {ApiResponse} from "../../../core/models/api-generic.model";
import {resolve} from "../../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {MessageContent, Messages} from "../models/message.model";

@Injectable({
  providedIn: 'root'
})
export class MessageService {

  constructor(private http: HttpClient) { }

  read(chat: string, first: number, last: number) {
    const fd = this.getFormData("read")
    fd.set("data", JSON.stringify({"chat": chat, "first": first, "last": last}))
    return this.http.post<ApiResponse<Messages>>(resolve(), fd)
  }

  send(chat: string, message: MessageContent) {
    let fd = this.getFormData("send")
    fd.set("data", JSON.stringify({"chat": chat, "text": message.text}))
    return this.http.post<ApiResponse<any>>(resolve(), fd)
  }

  private getFormData(method: string) {
    let fd = new FormData()
    fd.set("module", "messages")
    fd.set("method", method)
    return fd
  }
}
