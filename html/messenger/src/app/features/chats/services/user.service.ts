import { Injectable } from '@angular/core';
import {ApiResponse} from "../../../core/models/api-generic.model";
import {Chat} from "../models/chat.model";
import {resolve} from "../../../environments/environment";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(private http: HttpClient) { }

  exist(user: string) {
    const fd = this.getFormData("exist")
    fd.set("data", JSON.stringify({"user": user}))
    return this.http.post<ApiResponse<boolean>>(resolve(), fd)
  }

  private getFormData(method: string) {
    let fd = new FormData()
    fd.set("module", "user")
    fd.set("method", method)
    return fd
  }
}
