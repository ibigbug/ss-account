import { BrowserModule } from '@angular/platform-browser';
import { NgModule, APP_INITIALIZER } from '@angular/core';

import { MenubarModule, DataTableModule, ButtonModule, SharedModule, InputTextModule, MessagesModule } from "primeng/primeng";

import { AppComponent } from './app.component';
import { UsageTableComponent } from './usage-table.component';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    AppComponent,
    UsageTableComponent,
  ],
  imports: [
    BrowserModule,
    ButtonModule,
    MenubarModule,
    DataTableModule,
    SharedModule,
    HttpModule,
    InputTextModule,
    FormsModule,
    MessagesModule,
  ],
  providers: [
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
