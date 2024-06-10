import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdBannerItemComponent } from './ad-banner-item.component';

describe('AdBannerItemComponent', () => {
  let component: AdBannerItemComponent;
  let fixture: ComponentFixture<AdBannerItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AdBannerItemComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AdBannerItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
