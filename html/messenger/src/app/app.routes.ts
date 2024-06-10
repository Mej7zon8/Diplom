import { Routes } from '@angular/router';
import {
  AuthenticationRootComponent
} from "./features/authentication/components/authentication-root/authentication-root.component";
import {ChatsRootComponent} from "./features/chats/components/chats-root/chats-root.component";

export const routes: Routes = [
  {path: "sign-in", component: AuthenticationRootComponent},
  {path: "chats", component: ChatsRootComponent},
  {path: "**", redirectTo: "/sign-in"}
];
