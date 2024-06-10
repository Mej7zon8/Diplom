import {Component} from '@angular/core';
import {MatHint} from "@angular/material/form-field";

@Component({
  selector: 'app-preloader',
  standalone: true,
  imports: [
    MatHint
  ],
  templateUrl: './preloader.component.html',
  styleUrl: './preloader.component.scss'
})
export class PreloaderComponent {
  protected invisible = false
  protected hidden = false
  protected error: string = ""

  constructor() {
  }

  isShown() {
    return !this.invisible
  }

  show() {
    this.hidden = false
    this.invisible = false
  }

  hide() {
    this.invisible = true
    if (!this.hidden)
      setTimeout(() => {
        this.hidden = true
      }, 200)
  }

  // setError displays the provided error.
  setError(error: string) {
    this.error = error
  }

  protected get errorTitle() {
    if (this.error == "logged out")
      return "Вы вышли."
    else
      return "Упс, что-то пошло не так."
  }

  protected get errorSuggestion() {
    return ""
  }

  protected get rawError() {
    if (this.error == "logged out")
      return ""
    return this.error
  }
}
