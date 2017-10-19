import { Component, OnInit } from "@angular/core";
import { User } from "./types/types";
import { UserService } from "./services/user.service";

@Component({
    selector: 'usage-table',
    templateUrl: 'usage-table.component.html',
    providers: [UserService],
})
export class UsageTableComponent implements OnInit {
    loading: boolean;
    users: User[];
    cols: any[];

    constructor(private userService: UserService){}

    ngOnInit() {
        this.loading = true

        this.cols = [
            {field: 'username', header: 'Username'},
            {field: 'port', header: 'Port'},
            {field: 'backend', header: 'Backend'},
            {field: 'bytes_upload', header: 'Bytes Upload'},
            {field: 'bytes_download', header: 'Bytes Download'},
        ]
        this.userService.getUsers().subscribe(users => {
            this.loading = false
            this.users = users
        })
    }
}