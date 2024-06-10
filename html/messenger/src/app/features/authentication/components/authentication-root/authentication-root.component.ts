import {Component, ViewChild} from '@angular/core';
import {AuthenticationService} from "../../services/authentication.service";
import {PreloaderComponent} from "../../../../shared/components/preloader/preloader.component";
import {Router} from "@angular/router";
import {MatCard, MatCardActions, MatCardContent, MatCardHeader, MatCardTitle} from "@angular/material/card";
import {MatFormField, MatLabel} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {MatTab, MatTabGroup} from "@angular/material/tabs";
import {MatButton} from "@angular/material/button";
import {SnackbarErrorHandlerService} from "../../../../core/services/snackbar-error-handler.service";
import {tap} from "rxjs";

@Component({
  selector: 'app-authentication-root',
  standalone: true,
  imports: [
    PreloaderComponent,
    MatCard,
    MatCardHeader,
    MatCardTitle,
    MatCardContent,
    MatFormField,
    MatInput,
    MatLabel,
    MatTabGroup,
    MatTab,
    MatCardActions,
    MatButton
  ],
  templateUrl: './authentication-root.component.html',
  styleUrl: './authentication-root.component.scss'
})
export class AuthenticationRootComponent {
  @ViewChild('preloader', {static: true}) preloader!: PreloaderComponent;

  constructor(private service: AuthenticationService, private eh: SnackbarErrorHandlerService, private router: Router) {
  }

  ngOnInit() {
    this.service.check().subscribe(data => {
      // The request should only return one of the following codes: 204, 401.
      // Display an error on any other code.
      if (data.status != 204 && data.status != 401) {
        this.preloader.setError(data.error)
        return
      }
      // If the user is authenticated, redirect to the dashboard.
      if (data.status == 204) {
        this.router.navigateByUrl('/dashboard').then()
        return
      }
      // Close the preloader.
      // The user must authenticate manually.
      this.preloader.hide()
      this.disabled = false
    })
  }

  signUp(email: string, login: string, password: string, name: string) {
    this.disabled = true
    this.service.singUp(email, login, password, name).pipe(tap(this.eh.tap())).subscribe(data => {
      this.disabled = false
      if (data.status == 204) {
        this.router.navigateByUrl('/chats').then()
      }
    })
  }

  signIn(login: string, password: string) {
    this.disabled = true
    this.service.signIn(login, password).pipe(tap(this.eh.tap())).subscribe(data => {
      this.disabled = false
      if (data.status == 204) {
        this.router.navigateByUrl('/chats').then()
      }
    })
  }

  protected disabled = true
}
