import {Component, EventEmitter, Input, Output} from '@angular/core';
import {AdvertisementService} from "../../services/advertisement.service";
import {AsyncPipe} from "@angular/common";

@Component({
  selector: 'app-ad-banner-item',
  standalone: true,
  imports: [
    AsyncPipe
  ],
  templateUrl: './ad-banner-item.component.html',
  styleUrl: './ad-banner-item.component.scss'
})
export class AdBannerItemComponent {
  @Input() visible: boolean = false
  @Output() loaded = new EventEmitter<void>()

  protected ad$ = this.service.get()

  constructor(protected service: AdvertisementService) {
  }

  protected readonly window= window;
}
