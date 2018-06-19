import { Component, OnInit } from '@angular/core';
import {ApiService} from "../../services/api/api.service";
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";
import {Observable} from "rxjs/internal/Observable";

@Component({
  selector: 'app-balance',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.less']
})
export class BalanceComponent implements OnInit {
  public balance$: Observable<any>;
  public error: boolean = false;
  public complete: boolean = false;
  public balance: Array<Balance> = null;

  constructor(private api: ApiService) { }

  ngOnInit() {
    this.balance$ = this.api.getAccountBalance();
  }

  onSuccess(data) {
    this.balance = data.map(item => {
      item['color'] = this.amountColor(item.amount);
      return item;
    });

    this.error = false;
    this.complete = true;
  }

  onError(err) {
    this.error = true;
    this.complete = true;
    this.balance = null;
  }

  public amountColor(amount: number): string {
    let ret;

    if(amount > 0) {
      ret = 'text-success';
    } else if(amount < 0) {
      ret = 'text-danger';
    } else {
      ret = 'text-default';
    }

    return ret;
  }
}

export interface Balance {
  amount: number,
  label: string,
  color: string
}
