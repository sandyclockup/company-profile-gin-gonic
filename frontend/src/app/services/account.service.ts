/**
 * This file is part of the Sandy Andryanto Company Profile Website.
 *
 * @author     Sandy Andryanto <sandy.andryanto.dev@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from  '../../environments/environment';
import { HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AccountService {

  constructor(private http: HttpClient) {}

  public authRequestHeader(): any {
      let token = window.sessionStorage.getItem("token");
      return {
        "Content-Type": "application/json",
        "Authorization": "Bearer "+token
      }
  }

  profileDetail(): Observable<any>{
    return this.http.get(environment.backendURL+"/account/profile/detail", { headers: this.authRequestHeader() });
  }

  profileUpdate(data: any): Observable<any>{
    return this.http.post(environment.backendURL+"/account/profile/update", data, { headers: this.authRequestHeader() });
  }

  passwordUpdate(data: any): Observable<any>{
    return this.http.post(environment.backendURL+"/account/password", data, { headers: this.authRequestHeader() });
  }

  profileUpload(data: any): Observable<any>{
    let token = window.sessionStorage.getItem("token");
    let headers = new HttpHeaders({ "Authorization": "Bearer "+token});
    return this.http.post(environment.backendURL+"/account/upload", data, { headers });
  }

}
