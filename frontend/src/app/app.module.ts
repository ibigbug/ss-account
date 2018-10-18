import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { NgModule } from '@angular/core';

import { MenubarModule, ButtonModule, SharedModule, InputTextModule, MessagesModule } from "primeng/primeng";
import { TableModule } from "primeng/table";

import { AppComponent } from './app.component';
import { UsageTableComponent } from './usage-table.component';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { HumanReadableBytesPipe } from './human-readable-bytes.pipe';

@NgModule({
  declarations: [
    AppComponent,
    UsageTableComponent,
    HumanReadableBytesPipe,
  ],
  imports: [
    BrowserModule,
    ButtonModule,
    MenubarModule,
    TableModule,
    SharedModule,
    HttpModule,
    InputTextModule,
    FormsModule,
    MessagesModule,
    BrowserAnimationsModule,
  ],
  providers: [
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
