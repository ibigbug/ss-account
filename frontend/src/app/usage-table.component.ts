import { Component, OnInit } from "@angular/core";
import { User } from "./types/types";
import { UserService } from "./services/user.service";
import { MessageService } from "primeng/components/common/messageservice";

@Component({
    selector: 'usage-table',
    templateUrl: 'usage-table.component.html',
    providers: [UserService, MessageService],
})
export class UsageTableComponent implements OnInit {
    loading: boolean;
    users: User[];
    cols: any[];

    user: User; // current user

    constructor(
        private userService: UserService,
        private messageService: MessageService,
    ) { }

    ngOnInit() {
        this.loading = true
        this.user = new User()

        this.cols = [
            { field: 'username', header: 'Username' },
            { field: 'port', header: 'Port' },
            { field: 'backend', header: 'Backend' },
            { field: 'bytes_upload', header: 'Bytes Upload' },
            { field: 'bytes_download', header: 'Bytes Download' },
        ]
        this.userService.getUsers().subscribe(users => {
            this.loading = false
            this.users = users
        })
    }

    createUser() {
        this.loading = true
        this.userService.createUser(this.user).subscribe(u => {
            this.users.splice(0, 0, u)
        }, err => {
            this.messageService.add({
                severity: 'error',
                summary: 'User creation error: ',
                detail: err,
            })
        })
        this.loading = false
    }

    isValidForm(): boolean {
        return this.user.username && this.user.backend
    }
}