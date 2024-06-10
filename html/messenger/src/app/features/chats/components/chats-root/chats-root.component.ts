import {Component, EventEmitter, ViewChild} from '@angular/core';
import {ChatService} from "../../services/chat.service";
import {PreloaderComponent} from "../../../../shared/components/preloader/preloader.component";
import {exhaustMap, merge, Observable, tap, timer} from "rxjs";
import {SnackbarErrorHandlerService} from "../../../../core/services/snackbar-error-handler.service";
import {ApiResponse} from "../../../../core/models/api-generic.model";
import {Chat, ChatList} from "../../models/chat.model";
import {ChatListComponent} from "../chat-list/chat-list.component";
import {AsyncPipe} from "@angular/common";
import {ChatComponent} from "../chat/chat.component";
import {NewChatDialogComponent} from "../new-chat-dialog/new-chat-dialog.component";
import {MatDialog} from "@angular/material/dialog";
import {MatIconButton} from "@angular/material/button";
import {MatIcon} from "@angular/material/icon";
import {AuthenticationService} from "../../../authentication/services/authentication.service";
import {AdBannerComponent} from "../../../advertisement/components/ad-banner/ad-banner.component";

@Component({
  selector: 'app-chats-root',
  standalone: true,
  imports: [
    PreloaderComponent,
    ChatListComponent,
    AsyncPipe,
    ChatComponent,
    MatIconButton,
    MatIcon,
    AdBannerComponent
  ],
  templateUrl: './chats-root.component.html',
  styleUrl: './chats-root.component.scss'
})
export class ChatsRootComponent {
  @ViewChild('preloader', {static: true}) preloader!: PreloaderComponent;

  protected manualUpdates = new EventEmitter<void>()
  protected dataEmitter: Observable<ApiResponse<ChatList>>

  protected model: {
    selected_chat: Chat | null,
  } = {
    selected_chat: null,
  }

  constructor(private service: ChatService, private auth: AuthenticationService, private eh: SnackbarErrorHandlerService, private dialog: MatDialog) {
    this.dataEmitter = merge(timer(0, 4000), this.manualUpdates).pipe(
      exhaustMap(() => this.service.list()),
      tap(v => {
        if (v.error != "") {
          if (this.preloader.isShown()) // If the preloader is shown, display the response error there
            this.preloader.setError(v.error)
          else // Otherwise, display the error in the snackbar
            this.eh.show(v.error)
        } else
          this.preloader.hide()
      })
    )
  }

  protected createChat() {
    this.dialog.open(NewChatDialogComponent).afterClosed().subscribe(v => {
      this.manualUpdates.emit()
    })
  }

  protected leaveChat(chat: Chat) {
    this.service.leave(chat).subscribe(() => {
      this.manualUpdates.emit()
    })
  }

  protected exit() {
    this.auth.signOut()
  }

}
