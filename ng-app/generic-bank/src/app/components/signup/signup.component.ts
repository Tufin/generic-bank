import {Component} from '@angular/core';
import {ApiService} from "../../services/api/api.service";
import {Observable} from "rxjs/internal/Observable";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.less']
})
export class SignupComponent {
  public thanksLink = '';

  public model = {
    first: '',
    last: ''
  };

  public respond$: Observable<any>;

  public complete = null;
  public error = null;
  public data = null;

  constructor(private api: ApiService, private router: Router, private route: ActivatedRoute) {
    this.thanksLink = route.data['value'].link;
  }

  createAccount(first: string, last: string, $event) {
    $event.preventDefault();

    this.respond$ = this.api.createAccount(first, last);
  }

  public onSuccess(data) {
    const accountRegex = /'(\S+)'/g;

    const id = accountRegex.exec(data)[1];

    this.router.navigate([this.thanksLink, id]);
  }
}
