import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.less']
})
export class HomeComponent {
  public loginLink = '';
  public signupLink = '';

  constructor(private route: ActivatedRoute) {
    this.loginLink = route.data['value'].loginLink;
    this.signupLink = route.data['value'].signupLink;
  }
}
