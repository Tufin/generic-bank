import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, ChildActivationEnd, Router, RouterState} from "@angular/router";

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.less']
})
export class AdminComponent {
  public accountListLink = '';

  constructor(private route: ActivatedRoute) {
    this.accountListLink = route.data['value'].link;
    // console.log(route.data);
  }
}
