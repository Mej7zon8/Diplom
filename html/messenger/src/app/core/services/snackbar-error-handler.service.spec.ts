import { TestBed } from '@angular/core/testing';

import { SnackbarErrorHandlerService } from './snackbar-error-handler.service';

describe('SnackbarErrorHanderService', () => {
  let service: SnackbarErrorHandlerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SnackbarErrorHandlerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
