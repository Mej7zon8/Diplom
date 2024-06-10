import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChatsRootComponent } from './chats-root.component';

describe('ChatsRootComponent', () => {
  let component: ChatsRootComponent;
  let fixture: ComponentFixture<ChatsRootComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChatsRootComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ChatsRootComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
