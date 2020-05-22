export function getStatus(meta, validating?: any) {
    if (meta.active && validating) {
        return 'validating';
    }
    if (!meta.touched) {
        return null;
    }
    return meta.error ? 'error' : 'success';
}
