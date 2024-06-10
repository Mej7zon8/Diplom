import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AuthenticationRootComponent } from './authentication-root.component';

describe('AuthenticationRootComponent', () => {
  let component: AuthenticationRootComponent;
  let fixture: ComponentFixture<AuthenticationRootComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AuthenticationRootComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AuthenticationRootComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
