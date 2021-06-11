<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-table :items="ranks" :fields="rank_fields">
                    <template v-slot:cell(rank)="row">
                        <b-badge class="mr-2" :variant="rankStyle(row.item.rank)">{{ row.item.rank }}</b-badge>
                    </template>
                    <template v-slot:cell(shares)="row">
                        <b-badge v-for="(s, index) in row.item.shares" :key="index" class="mr-2" variant="light">{{ s.subject }}</b-badge>
                    </template>
                </b-table>
            </b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';

export default {
        name: 'Rank',
        components: {},
        data() {
            return {
                rank_fields: [
                    {key: 'rank', label: '排名'},
                    {key: 'name', label: '用户'},
                    {key: 'credit', label: '贡献'},
                    {key: 'shares', label: '相关分享'},
                ],
                ranks: [],
            };
        },
        computed: {
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            rankStyle(rank) {
                switch (rank) {
                case 1: return 'warning';
                case 2: return 'success';
                case 3: return 'primary';
                default:
                }
            },
            reload() {
                axios.get('/api/credits/rank/').then(response => {
                    this.ranks = response.data;
                }).catch(error => {this.ErrorBox(error)});
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>

<style scoped>

</style>