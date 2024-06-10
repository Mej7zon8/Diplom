import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {ApiResponse} from "../../../core/models/api-generic.model";
import {resolve} from "../../../environments/environment";
import {Ad} from "../models/ad.model";

@Injectable({
  providedIn: 'root'
})
export class AdvertisementService {
  constructor(private http: HttpClient) { }

  get() {
    const fd = this.getFormData("get")
    return this.http.post<ApiResponse<Ad>>(resolve(), fd)
  }

  private getFormData(method: string) {
    let fd = new FormData()
    fd.set("module", "advertisement")
    fd.set("method", method)
    return fd
  }
}
