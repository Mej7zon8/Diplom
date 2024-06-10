import {Component, ViewChild, ViewContainerRef} from '@angular/core';
import {AdBannerItemComponent} from "../ad-banner-item/ad-banner-item.component";

@Component({
  selector: 'app-ad-banner',
  standalone: true,
  imports: [
    AdBannerItemComponent
  ],
  templateUrl: './ad-banner.component.html',
  styleUrl: './ad-banner.component.scss'
})
export class AdBannerComponent {
  @ViewChild("content", {read: ViewContainerRef, static: true}) container!: ViewContainerRef

  ngOnInit() {
    this.insertAd()
    setInterval(() => {
      this.insertAd()

      // Remove the old element (if it exists)
      // if (this.container.length > 1) {
      //   this.container.remove(0);
      // }

      // newEl.instance.data = {
      //   image: "https://via.placeholder.com/728x90.png?text=Visit+Who+You+Gonna+Call",
      //   config: {
      //     url: "https://www.youtube.com/watch?v=Fe93CLbHjxQ"
      //   }
      // }
    }, 40000)
  }

  private insertAd() {
    const newEl = this.container.createComponent(AdBannerItemComponent)
    newEl.location.nativeElement.style.display = "contents"

    newEl.instance.loaded.subscribe(() => {
      // Queue animation (doesn't work without setTimeout)
      setTimeout(() => {
        newEl.instance.visible = true
      }, 0)

      // Remove the old element (if it exists), after the animation is done
      setTimeout(() => {
        console.debug(this.container.length)
        if (this.container.length > 1) {
          this.container.remove(0);
        }
      }, 500)
    })
  }
}
