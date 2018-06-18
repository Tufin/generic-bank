import {AfterContentChecked, AfterContentInit, Component, Input, OnInit} from '@angular/core';
import { Location } from '@angular/common';
import {ApiService, timeResponse} from "../../services/api/api.service";
import {Observable} from "rxjs/internal/Observable";
import animations from "../../animations/animations";
import {isNullOrUndefined} from "util";

@Component({
  selector: 'app-main-panel',
  templateUrl: './main-panel.component.html',
  styleUrls: ['./main-panel.component.less'],
  animations: [ animations.expandTransition]
})
export class MainPanelComponent {
  private _showTime: boolean = false;

  public times$: Observable<any>;
  public times: Array<timeResponse> = null;
  public timesError: boolean = false;
  public timesLoading: boolean = false;

  @Input() title = 'This is a title';
  @Input() showBack = false;
  @Input() set showTime(value: boolean) {
    this._showTime = value;

    if(value) {
      this.timesError = false;
      this.timesLoading = false;
      this.times$ = this.api.getTime('Asia/Jerusalem', 'Asia/Tokyo', 'Europe/London');
    }
  };

  get showTime(): boolean {
    return this._showTime;
  }

  constructor(private api: ApiService, private location: Location) { }

  goBack() {
    this.location.back();
  }

  onError(err) {
    this.timesError = err;
    this.timesLoading = false;
    this.times = null;
  }

  onSuccess(data) {
    if(isNullOrUndefined(data)) {
      return;
    }

    this.times = data;
    this.timesLoading = false;
    this.timesError = false;
  }
}
