import {Component, Input, ViewChild} from '@angular/core';
import {MatMenuModule, MatMenuTrigger} from "@angular/material/menu";

@Component({
  selector: 'app-context-menu',
  standalone: true,
  imports: [
    MatMenuModule
  ],
  templateUrl: './context-menu.component.html',
  styleUrl: './context-menu.component.scss'
})
export class ContextMenuComponent {
  @ViewChild(MatMenuTrigger) contextMenu!: MatMenuTrigger;

  protected position: {x: string, y: string} = {x: "0px", y: "0px"}
  @Input() items!: contextMenuItem[]

  // OpenContextMenu is the primary way to open the context menu.
  OpenContextMenu(event: MouseEvent) {
    event.preventDefault();

    this.position.x = event.clientX + 'px'
    this.position.y = event.clientY + 'px'
    this.contextMenu.openMenu()
  }

  OpenWithItems(event: MouseEvent, items: contextMenuItem[]) {
    this.items = items
    this.OpenContextMenu(event)
  }

  protected trackByItem(i: number) {
    return i
  }
}

export interface contextMenuItem {
  text: string
  disabled?: boolean
  action: () => void
}
