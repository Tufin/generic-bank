import {Directive, EventEmitter, HostBinding, Input, OnDestroy, Output} from '@angular/core';
import {Observable} from "rxjs/internal/Observable";
import {isNullOrUndefined} from "util";
import {Subscription} from "rxjs/internal/Subscription";
import {skip} from "rxjs/operators";

@Directive({
  selector: '[asyncProgress]'
})
export class AsyncProgressDirective implements OnDestroy{
  private _subscription: Subscription;

  @Output() data = new EventEmitter();
  @Output() error = new EventEmitter();
  @Output() complete = new EventEmitter();

  @Input() skipFirstValues: number = 0;

  @Input('asyncProgress') set obs$(obs: Observable<any>) {
    if(!isNullOrUndefined(obs)) {
      // this.emitChanges(null, null, false);

      this._subscription = obs.pipe(skip(this.skipFirstValues)).subscribe(
        (data) => {
          // console.log('data', data);
          this.data.emit(data);
          // this.emitChanges(data, null, false);
        },
        err => {
          // console.log('err', err);
          this.error.emit(err);

          // this.emitChanges(null, err, true);
        },
        () => {
          // console.log('completed');
          this.complete.emit(true);

          // this.emitChanges(null, null, true)
        }
      );
    }
  };

  constructor() { }

  ngOnDestroy() {
    if(!isNullOrUndefined(this._subscription)) {
      this._subscription.unsubscribe();
    }
  }

}
