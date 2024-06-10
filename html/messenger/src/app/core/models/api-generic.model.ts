// ApiResponse represents every api response.
import {finalize, MonoTypeOperatorFunction, Observable, of, switchMap, tap} from "rxjs";

export interface ApiResponse<T> {
  // data is the response data.
  // data is only valid when error is empty.
  data: T | null
  // code is the status code of the response.
  status: number
  error: string
}

export function UndefinedApiResponse<T>(): ApiResponse<T> | undefined {
  return undefined
}

export class RequestStatus {
  private is_progress = false

  onSend() {
    this.is_progress = true
  }

  onReceive() {
    this.is_progress = false
  }

  get status() {
    return this.is_progress
  }

  set value(v: boolean) {
    this.is_progress = v
  }

  with<T>(o: Observable<T>): Observable<T> {
    return of(0).pipe(
      tap(() => this.onSend()),
      switchMap(() => o),
    ).pipe(
      finalize(() => this.onReceive())
    )
  }
}


// RequestStatusV2 is an implementation of request status tracker.
// RequesterV2 distinguishes between get- and set-requests.
export class RequestStatusV2 {
  private state: "in_progress_set" | "after_set" | "idle" = "idle"

  onSend(type: "get"|"set") {
    if (type == "set")
      this.state = "in_progress_set"
  }

  onReceive(type: "get" | "set") {
    if (type == "set")
      this.state = "after_set"
    else if (this.state == "after_set")
      this.state = "idle"
  }

  get status() {
    return this.state != "idle"
  }

  with<T>(o: Observable<T>, type: "get" | "set"): Observable<T> {
    return of(0).pipe(
      tap(() => this.onSend(type)),
      switchMap(() => o),
    ).pipe(
      finalize(() => this.onReceive(type))
    )
  }
}
