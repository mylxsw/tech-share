<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2" no-body>
                <b-card-body>
                    <b-form>
                        <div v-if="shareDetail.share != null">
                            <b-form-group label="主题">
                                <b-form-input v-model="shareDetail.share.subject" type="text" readonly></b-form-input>
                            </b-form-group>
                            <b-form-group label="主题类型">
                                <b-form-input v-model="shareDetail.share.subject_type" type="text" readonly></b-form-input>
                            </b-form-group>
                            <b-form-group label="内容简介">
                                <b-form-textarea v-model="shareDetail.share.description" rows="4" readonly/>
                            </b-form-group>
                            <b-form-group label="分享人" >
                                <b-form-input v-model="shareDetail.share.share_user" type="text" readonly></b-form-input>
                            </b-form-group>
                        </div>

                        <div v-if="(shareDetail.join_users !== null && shareDetail.join_users.length > 0) || (shareDetail.like_users !== null && shareDetail.like_users.length > 0)" class="mt-3">
                            <div v-if="shareDetail.like_users !== null && shareDetail.like_users.length > 0">
                                对分享感兴趣的人<br />
                                <b-badge class="mr-2" variant="info" v-for="(item, index) in shareDetail.like_users" :key="index">{{ item.name }}</b-badge>
                            </div>
                            <div v-if="shareDetail.join_users !== null && shareDetail.join_users.length > 0">
                                报名参加分享的人<br />
                                <b-badge class="mr-2" variant="primary" v-for="(item, index) in shareDetail.join_users" :key="index">{{ item.name }}</b-badge>
                            </div>
                        </div>

                        <div v-if="shareDetail.plan != null" class="mt-3">
                            <b-form-group label="分享时间" readonly v-if="shareDetail.plan.share_at != ''">
                                <date-time :value="shareDetail.plan.share_at"></date-time>
                            </b-form-group>
                            <b-form-group label="分享地点" readonly v-if="shareDetail.plan.share_room != ''">
                                <b-form-input v-model="shareDetail.plan.share_room" type="text" readonly></b-form-input>
                            </b-form-group>
                            <b-form-group label="预计时长" readonly v-if="shareDetail.plan.plan_duration > 0">
                                <b-form-input v-model="shareDetail.plan.plan_duration" type="text" readonly></b-form-input>
                            </b-form-group>
                            <b-form-group label="实际时长" readonly v-if="shareDetail.plan.real_duration > 0">
                                <b-form-input v-model="shareDetail.plan.real_duration" type="text" readonly></b-form-input>
                            </b-form-group>
                            <b-form-group label="备注" v-if="shareDetail.plan.note != ''">
                                <b-form-textarea v-model="shareDetail.plan.note" rows="4" readonly />
                            </b-form-group>
                        </div>

                        <div v-if="shareDetail.attachments != null && shareDetail.attachments.length > 0 && shareDetail.plan != null" class="mt-3">
                            <b-form-group label="附件">
                                <li v-for="(atta, index) in shareDetail.attachments" :key="index">
                                    <a :href="'/storage/' + atta.atta_path" target="_blank">{{ atta.name }}</a>
                                </li>
                            </b-form-group>
                        </div>
                    </b-form>
                </b-card-body>
                <b-card-footer v-if="shareDetail.share != null && shareDetail.share.create_user_id != this.$store.getters.user.id && shareDetail.share.status != 3" >
                    <b-button size="sm" variant="success" v-if="!shareDetail.user_like && canLike(shareDetail.share)" @click="iLikeIt(shareDetail.share.id, true)">感兴趣</b-button>
                    <b-button size="sm" variant="danger" v-if="shareDetail.user_like && canLike(shareDetail.share)" @click="iLikeIt(shareDetail.share.id, false)">不感兴趣</b-button>
                    <b-button size="sm" variant="success" v-if="shareDetail.share.status === 2 && !shareDetail.user_join" @click="iJoinIt(shareDetail.share.id, true)">参加</b-button>
                    <b-button size="sm" variant="danger" v-if="shareDetail.share.status === 2 && shareDetail.user_join" @click="iJoinIt(shareDetail.share.id, false)">不参加了</b-button>
                </b-card-footer>
            </b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import moment from 'moment';

export default {
        name: 'Share',
        components: {},
        data() {
            return {
                shareDetail: {share: {}, plan: {}, join_users: [], like_users: []},
            };
        },
        computed: {
        },
        watch: {
            '$route': 'reload',
        },
        methods: {
            // 时间格式化
            dateFormat(d) {
                return moment(d).format('YYYY-MM-DD HH:mm:ss');
            },
            
            // 是否可以喜欢/参加
            canLike(share) {
                return share.status === 1 && this.$store.getters.user.id != share.create_user_id;
            },
            canNotLike(share) {
                return share.status === 1 && this.$store.getters.user.id != share.create_user_id;
            },
            // 喜欢or参加分享
            iLikeIt(shareId, like, success) {
                success = success || function() {};
                axios[like ? 'post':'delete']('/api/shares/' + shareId + '/like/').then(() => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            iJoinIt(shareId, join, success) {
                success = success || function() {};
                axios[join ? 'post':'delete']('/api/shares/' + shareId + '/join/').then(() => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
        
            // 页面刷新
            reload() {
                axios.get('/api/shares/' + this.$route.query.id + '/').then(response => {
                    this.shareDetail = response.data;
                }).catch(error => {this.ErrorBox(error)});
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>
