import { Injectable } from '@angular/core';
import {HttpClient,} from "@angular/common/http";
import {Observable} from "rxjs/internal/Observable";
import {map, take, catchError} from 'rxjs/operators';
import { interval } from 'rxjs';
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";
import { environment } from './../../../environments/environment';

@Injectable()
export class ApiService {
  // private baseURL: string;
  private readonly PREFIX = environment.production ? '' : '/api';
  private readonly TIME_URL = '/time';
  private readonly CREATE_ACCOUNT_URL = '/accounts/';
  private readonly ACCOUNTS_LIST_URL = '/boa/admin/accounts';
  private readonly BALANCE_URL = '/balance';

  constructor(private http: HttpClient) {
  }


  getUTCTime(zone: string) {
    return this.httpRequest<timeResponse>(`${this.TIME_URL}?zone=${zone}`, 'GET').pipe(
      map(item => {
        item.formatted = item.formatted.replace(/-/g, '/');
        item.date = new Date(item.formatted).toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });
        return item;
      })
    );
  }

  getAccountsList() {
    return this.httpRequest<accountListResponse>(this.ACCOUNTS_LIST_URL, 'GET');
  }

  createAccount(first: string, last: string): Observable<any> {
    const accountId = `${first}:${last}:${this.getRandomInt(10000,99999)}`;

    return this.httpRequest<string>(`${this.CREATE_ACCOUNT_URL}${accountId}`, 'POST', {responseType: 'text'});
  }

  getAccountBalance() {
    return this.httpRequest<Array<balanceResponse>>(this.BALANCE_URL, "GET");
  }

  getTime(...zones): Observable<any> {
    let ret = new BehaviorSubject(null);
    let requests = [];

    let sub$ = interval(1000).pipe(
      take(zones.length),
      map((n) => {
        console.log(zones[n]);
        return this.getUTCTime(zones[n]);
      })
    ).subscribe((item) => {
      item.subscribe(time => {
        requests.push(time);

        if(requests.length === zones.length) {
          ret.next(requests);
          return ret.complete();
        }
      }, err => {
        ret.error(err);
        sub$.unsubscribe();
        return ret.complete();
      });
    });

    // console.log(requests);
    return ret.asObservable();
  }

  private httpRequest<T>(url: string, method: string, options: any = {}): Observable<any> {
    return this.http.request<T>(method, url, options).pipe(
      catchError(err => {
        // console.log('ERROR HANDLER');
        throw new Error(err)
      })
    );
  }

  private getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min)) + min;
  }
}

export interface timeResponse {
  countryName: string,
  formatted: string,
  status: string
}

export interface accountListResponse {
  accounts: Array<string>
}

export interface balanceResponse {
  label: string,
  amount: number
}
