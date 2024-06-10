import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {ApiResponse} from "../../../core/models/api-generic.model";
import {resolve} from "../../../environments/environment";
import {tap} from "rxjs";
import {Router} from "@angular/router";

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  public signedOut = false

  constructor(private http: HttpClient, private router: Router) { }

  // check checks whether the user is signed in
  check() {
    return this.http.post<ApiResponse<any>>(resolve(), this.getFormData("check"))
  }

  singUp(email: string, username: string, password: string, name: string) {
    const fd = this.getFormData("sign-up")
    fd.set("data", JSON.stringify({"email": email, "username": username, "password": password, "name": name}))
    return this.http.post<ApiResponse<any>>(resolve(), fd).pipe(tap(v => {
      // Save the username in the local storage
      if (v.error == "")
        localStorage.setItem("username", username)
    }))
  }

  signIn(username: string, password: string) {
    const fd = this.getFormData("sign-in")
    fd.set("data", JSON.stringify({"username": username, "password": password}))
    return this.http.post<ApiResponse<any>>(resolve(), fd).pipe(tap(v => {
      // Save the username in the local storage
      if (v.error == "")
        localStorage.setItem("username", username)
    }))
  }

  signOut() {
    this.http.post<ApiResponse<any>>(resolve(), this.getFormData("sign-out")).pipe(tap(v => {
      if (v.error == "") {
        localStorage.removeItem("username")
        this.username = ""
        this.signedOut = true
        this.router.navigateByUrl("/sign-in").then()
      }
    })).subscribe()
  }

  private getFormData(method: string) {
    let fd = new FormData()
    fd.set("module", "authentication")
    fd.set("method", method)
    return fd
  }

  private username: string = ""

  getUsername() {
    if (this.username != "")
      return this.username
    const username = localStorage.getItem("username")
    if (username)
      this.username = username
    return username
  }
}
