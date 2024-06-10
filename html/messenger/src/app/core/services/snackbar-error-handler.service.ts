import {Injectable} from '@angular/core';
import {MatSnackBar} from "@angular/material/snack-bar";
import {ApiResponse} from "../models/api-generic.model";
import {titleCaseWord} from "../tools/tools";

@Injectable({
  providedIn: 'root'
})
export class SnackbarErrorHandlerService {

  constructor(private snackBar: MatSnackBar) {
  }

  tap<T>(): (v: ApiResponse<T>) => void {
    return v => {
      if (v.error != "")
        this.show(v.error)
    }
  }

  show(error: string, text = 'OK') {
    const ref = this.snackBar.open(this.transformError(error), text)
    // Create a multicasted observable:
    // const mc = new ReplaySubject<void>(1)
    // ref.afterDismissed().subscribe(() => mc.next())

    // if (pageReload)
    //   mc.subscribe(() => window.location.reload())
    return ref
  }

  transformError(error: string) {
    if (error == "bad request")
      error = `Что-то пошло не так`
    if (error == "invalid credentials")
      error = `Неверный логин или пароль`
    if (error == "user already exists")
      error = `Пользователь уже существует`
    if (error == "email is invalid")
      error = `Почта недопустима`
    if (error == "username is invalid")
      error = `Имя пользователя недопустимо`
    if (error == "password is invalid")
      error = `Пароль недопустим`
    if (error == "name is invalid")
      error = `Имя недопустимо`
    if (error == "user not found")
      error = `Пользователь не найден`
    if (error == "chat already exists")
      error = `Чат уже существует`
    if (error == "cannot create chat with self")
      error = `Нельзя создать чат с самим собой`
    return titleCaseWord(error)
  }
}
