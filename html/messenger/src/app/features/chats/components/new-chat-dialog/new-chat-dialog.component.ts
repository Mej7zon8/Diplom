import {Component} from '@angular/core';
import {
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";
import {MatButton} from "@angular/material/button";
import {MatFormField, MatLabel} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {FormsModule} from "@angular/forms";
import {ChatService} from "../../services/chat.service";
import {SnackbarErrorHandlerService} from "../../../../core/services/snackbar-error-handler.service";
import {
  MatChip,
  MatChipGrid,
  MatChipInput,
  MatChipInputEvent,
  MatChipRemove,
  MatChipRow
} from "@angular/material/chips";
import {MatIcon} from "@angular/material/icon";
import {COMMA, ENTER} from "@angular/cdk/keycodes";
import {UserService} from "../../services/user.service";
import {tap} from "rxjs";

@Component({
  selector: 'app-new-chat-dialog',
  standalone: true,
  imports: [
    MatDialogTitle,
    MatDialogContent,
    MatDialogActions,
    MatButton,
    MatDialogClose,
    MatFormField,
    MatInput,
    FormsModule,
    MatLabel,
    MatChipRow,
    MatChipGrid,
    MatIcon,
    MatChipInput,
    MatChip,
    MatChipRemove
  ],
  templateUrl: './new-chat-dialog.component.html',
  styleUrl: './new-chat-dialog.component.scss'
})
export class NewChatDialogComponent {
  protected name = ""
  protected selected: string[] = []
  readonly separatorKeysCodes = [ENTER, COMMA] as const;
  protected disabled = false

  constructor(public dialogRef: MatDialogRef<NewChatDialogComponent>, private chatService: ChatService, private userService: UserService, private eh: SnackbarErrorHandlerService) {
    this.dialogRef.updateSize('400px')
  }

  protected remove(participant: string) {
    this.selected = this.selected.filter(v => v != participant)
  }

  protected add(event: MatChipInputEvent) {
    const value = event.value.trim();
    // Ignore events with en empty value
    if (!value) return;

    this.disabled = true
    // Check if the user exists
    this.userService.exist(event.value).pipe(tap(this.eh.tap())).subscribe(v => {
      this.disabled = false
      if (v.data) {
        this.selected.push(event.value)
        event.chipInput!.clear();
      } else
        this.eh.show(`Пользователь не найден: ${event.value}`)
    })
  }

  protected createChat() {
    this.chatService.create(this.selected, this.name).pipe(tap(this.eh.tap())).subscribe(v => {
      if (v.error == "")
        this.dialogRef.close()
    })
  }
}
