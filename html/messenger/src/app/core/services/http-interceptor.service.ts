import {Injectable} from '@angular/core';
import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpResponse
} from "@angular/common/http";
import {catchError, Observable, of, tap} from "rxjs";
import {ApiResponse} from "../models/api-generic.model";
import {Router} from "@angular/router";

@Injectable({
  providedIn: 'root'
})
export class HttpInterceptorService implements HttpInterceptor {

  constructor(private router: Router) {
  }

  intercept<T>(request: HttpRequest<ApiResponse<T>>, next: HttpHandler): Observable<HttpEvent<ApiResponse<T>>> {
    // enable cookie:
    request = request.clone({
      withCredentials: true
    });

    return next.handle(request)
      .pipe(
        // Catch the Unauthorized error for any request and redirect to the sign-in page.
        tap(event => {
          if (event instanceof HttpResponse && event.body && event.body.status && event.body.status == 401) {
            this.router.navigateByUrl('/sign-in').then()
          }
        }),
        catchError((error: HttpErrorResponse) => {
          console.error('Error from error interceptor', error);

          let errorMsg = '';
          if (error.error instanceof ErrorEvent) {
            // Client-side or network error occurred
            errorMsg = `${error.error.message}`;
          } else {
            // The backend returned an unsuccessful response code
            errorMsg = `${error.message}`;
          }

          const response = new HttpResponse({
            status: 200,
            body: {
              data: null,
              status: error.status,
              error: errorMsg
            }
          })
          return of(response)
        })
      );
  }

}
