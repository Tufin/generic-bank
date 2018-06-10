import {
  AfterContentChecked,
  AfterViewChecked,
  AfterViewInit,
  Directive,
  ElementRef,
  EventEmitter,
  HostBinding,
  HostListener,
  Output
} from '@angular/core';

@Directive({
  selector: '[exposeHeight]'
})
export class ExposeHeightDirective implements AfterContentChecked {
  @HostListener('window:resize', ['$event'])
  private onResize() {
    if(this.elem.nativeElement.children[1] && this.elem.nativeElement.children[1].children[0]) {
      this.elemHeight.emit(this.elem.nativeElement.children[1].children[0].clientHeight);
    }
  };

  @Output('exposeHeight')
  public elemHeight = new EventEmitter();

  constructor(private elem: ElementRef) {

  }

  ngAfterContentChecked() {
    this.onResize();
  }
}