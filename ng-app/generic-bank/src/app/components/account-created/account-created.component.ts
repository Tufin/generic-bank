import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";

@Component({
  selector: 'app-account-created',
  templateUrl: './account-created.component.html',
  styleUrls: ['./account-created.component.less']
})
export class AccountCreatedComponent {
  public id: string;
  public balanceLink = '';

  constructor(private route: ActivatedRoute) {
    this.id = route.params['value'].id;
    this.balanceLink = route.data['value'].link;
  }
}
