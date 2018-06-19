import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.less']
})
export class LoginComponent {
  public loginLink = '';

  public model = {
    first: '',
    last: ''
  };

  constructor(private router: Router, private route: ActivatedRoute) {
    this.loginLink = route.data['value'].link;
  }

  loginAccount(first, last, $event) {
    $event.preventDefault();

    // console.log(this.model, first, last);

    // if(this.model.first === 'admin' && this.model.last === 'admin') {
    this.router.navigate([this.loginLink]);
    // }
  }

}
