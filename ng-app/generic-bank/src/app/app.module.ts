import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import {RouterModule, Routes} from "@angular/router";
import { MainPanelComponent } from './components/main-panel/main-panel.component';
import { LoginComponent } from './components/login/login.component';
import { SignupComponent } from './components/signup/signup.component';
import { HomeComponent } from './components/home/home.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import { ExposeHeightDirective } from './directives/expose-height.directive';
import { AdminComponent } from './components/admin/admin.component';
import {ApiService} from "./services/api/api.service";
import {HttpClientModule} from "@angular/common/http";
import { AsyncProgressDirective } from './directives/async-progress.directive';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { ProgressCardComponent } from './components/progress-card/progress-card.component';
import { AccountCreatedComponent } from './components/account-created/account-created.component';
import { AccountsListComponent } from './components/accounts-list/accounts-list.component';

const appRoutes: Routes = [
  {
    path: 'home',
    component: HomeComponent,
    data: {
      title: 'Welcome to Generic Bank’s Customer Portal',
      showBack: false,
      // showTime: true
    }
  },
  {
    path: 'signup',
    component: SignupComponent,
    data: {
      title: 'Open an account',
      showBack: true
    }
  },
  {
    path: 'thanks/:id',
    component: AccountCreatedComponent,
    data: {
      title: 'Thank you for choosing us as your bank',
      showBack: false
    }
  },
  {
    path: 'login',
    component: LoginComponent,
    data: {
      title: 'Login to your account',
      showBack: true
    }
  },
  {
    path: 'admin',
    component: AdminComponent,
    data: {
      title: '',
      showBack: true,
      showTime: true
    }
  },
  {
    path: 'admin/accounts',
    component: AccountsListComponent,
    data: {
      title: 'New customer accounts',
      showBack: true,
      // showTime: true
    }
  },
  { path: '',
    redirectTo: '/home',
    pathMatch: 'full'
  },
  { path: '**',
    redirectTo: '/home',
    pathMatch: 'full' }
];

@NgModule({
  declarations: [
    AppComponent,
    MainPanelComponent,
    LoginComponent,
    SignupComponent,
    HomeComponent,
    ExposeHeightDirective,
    AdminComponent,
    AsyncProgressDirective,
    ProgressCardComponent,
    AccountCreatedComponent,
    AccountsListComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forRoot(appRoutes, {useHash: true})
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }
