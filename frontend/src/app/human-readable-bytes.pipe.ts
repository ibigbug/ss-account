import { Pipe, PipeTransform } from "@angular/core"


@Pipe({name: 'human_readable'})
export class HumanReadableBytesPipe implements PipeTransform {
    transform(value: number): string {
        if (value < 1) {
            return value + ' B'
        }
        let i = Math.floor(Math.log(value) / Math.log(1024))
        let sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

        return (value / Math.pow(1024, i)).toFixed(2) + ' ' + sizes[i]
    }
}