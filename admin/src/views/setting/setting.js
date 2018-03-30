import request from 'axios';
import _ from 'lodash';

import {
  SETTINGS,
} from '../../urls';
import {
  getDate,
  diff,
} from '../../helpers/util';

export default {
  data() {
    return {
      form: null,
      mode: -1,
      settings: null,
      modified: null,
    };
  },
  methods: {
    edit(id) {
      this.form = _.clone(_.find(this.settings, item => item.id === id));
      this.mode = 1;
    },
    add() {
      this.form = {};
      this.mode = 1;
    },
    async loadSettings() {
      const close = this.$loading();
      try {
        const res = await request.get(SETTINGS, {
          params: {
            'cache-control': 'no-cache',
          },
        });
        this.settings = _.map(res.data.list, (item) => {
          // eslint-disable-next-line
          item.data = JSON.stringify(item.data, null, 2);
          item.date = getDate(item.updatedAt).substring(0, 16);
          return item;
        });
        this.mode = 0;
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
    async submit() {
      const {
        form,
      } = this;
      const {
        name,
        disabled,
        data,
        description,
        id,
      } = form;
      if (!name || !data || !description) {
        this.$error('name, data and description can\'t be null');
        return;
      }
      const close = this.$loading();
      try {
        if (id) {
          const found = _.find(this.settings, item => item.id === id);
          const keys = ['name', 'disabled', 'data', 'description'];
          const updateData = diff(found, form, keys);
          await request.patch(`${SETTINGS}/${id}`, updateData);
        } else {
          await request.post(SETTINGS, {
            name,
            disabled,
            data: JSON.parse(data),
            description,
          });
        }
        await this.loadSettings();
      } catch (err) {
        this.$error(err);
      } finally {
        close();
      }
    },
  },
  beforeMount() {
    this.loadSettings();
  },
};
